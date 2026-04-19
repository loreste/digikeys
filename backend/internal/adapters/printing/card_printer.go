package printing

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/digikeys/backend/internal/domain"
)

// PrintBatch tracks a group of cards queued for printing.
type PrintBatch struct {
	ID        string
	CardIDs   []string
	Status    string // pending, processing, completed, failed
	CreatedAt time.Time
}

// CardPrinter implements ports.PrintingService with in-memory batch tracking.
// In production, this would integrate with a card printing hardware API.
type CardPrinter struct {
	mu      sync.RWMutex
	batches map[string]*PrintBatch
	// cardData stores card+citizen pairs keyed by card ID for print file generation.
	cardData map[string]*printEntry
}

type printEntry struct {
	Card    *domain.Card    `json:"card"`
	Citizen *domain.Citizen `json:"citizen"`
}

// NewCardPrinter creates a new in-memory card printer service.
func NewCardPrinter() *CardPrinter {
	return &CardPrinter{
		batches:  make(map[string]*PrintBatch),
		cardData: make(map[string]*printEntry),
	}
}

// QueueForPrinting adds a card to the default pending batch, creating one if needed.
func (p *CardPrinter) QueueForPrinting(_ context.Context, card *domain.Card, citizen *domain.Citizen) error {
	if card == nil || citizen == nil {
		return fmt.Errorf("card and citizen are required")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Store the card data for later print file generation.
	p.cardData[card.ID] = &printEntry{Card: card, Citizen: citizen}

	// Find or create a pending batch.
	var batch *PrintBatch
	for _, b := range p.batches {
		if b.Status == "pending" {
			batch = b
			break
		}
	}
	if batch == nil {
		batch = &PrintBatch{
			ID:        uuid.New().String(),
			Status:    "pending",
			CreatedAt: time.Now(),
		}
		p.batches[batch.ID] = batch
	}

	batch.CardIDs = append(batch.CardIDs, card.ID)
	slog.Info("card queued for printing", "cardId", card.ID, "batchId", batch.ID, "batchSize", len(batch.CardIDs))
	return nil
}

// GetBatchStatus returns the status of a print batch.
func (p *CardPrinter) GetBatchStatus(_ context.Context, batchID string) (string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	batch, ok := p.batches[batchID]
	if !ok {
		return "", fmt.Errorf("batch %s not found", batchID)
	}
	return batch.Status, nil
}

// CreatePrintBatch creates a new print batch from a list of card IDs and returns the batch ID.
func (p *CardPrinter) CreatePrintBatch(cardIDs []string) (string, error) {
	if len(cardIDs) == 0 {
		return "", fmt.Errorf("card IDs must not be empty")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	batch := &PrintBatch{
		ID:        uuid.New().String(),
		CardIDs:   cardIDs,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	p.batches[batch.ID] = batch

	slog.Info("print batch created", "batchId", batch.ID, "cardCount", len(cardIDs))
	return batch.ID, nil
}

// GeneratePrintFile generates a ZIP file containing JSON data for each card in the batch.
func (p *CardPrinter) GeneratePrintFile(batchID string) ([]byte, error) {
	p.mu.RLock()
	batch, ok := p.batches[batchID]
	if !ok {
		p.mu.RUnlock()
		return nil, fmt.Errorf("batch %s not found", batchID)
	}
	// Copy data we need under read lock.
	cardIDs := make([]string, len(batch.CardIDs))
	copy(cardIDs, batch.CardIDs)
	p.mu.RUnlock()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for _, cardID := range cardIDs {
		p.mu.RLock()
		entry, exists := p.cardData[cardID]
		p.mu.RUnlock()

		var data []byte
		var err error
		if exists {
			data, err = json.MarshalIndent(entry, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("marshal card %s: %w", cardID, err)
			}
		} else {
			// Card data not found, write a placeholder.
			data = []byte(fmt.Sprintf(`{"cardId": "%s", "error": "card data not loaded"}`, cardID))
		}

		fw, err := zw.Create(fmt.Sprintf("card_%s.json", cardID))
		if err != nil {
			return nil, fmt.Errorf("create zip entry: %w", err)
		}
		if _, err := fw.Write(data); err != nil {
			return nil, fmt.Errorf("write zip entry: %w", err)
		}
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("close zip: %w", err)
	}

	// Update batch status.
	p.mu.Lock()
	if b, ok := p.batches[batchID]; ok {
		b.Status = "processing"
	}
	p.mu.Unlock()

	return buf.Bytes(), nil
}
