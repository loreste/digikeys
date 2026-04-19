import React, { useState } from 'react';
import { View, Text, TouchableOpacity, ActivityIndicator, StyleSheet } from 'react-native';
import { useEnrollmentStore } from '../stores/enrollment-store';

export default function SyncScreen() {
  const enrollments = useEnrollmentStore((s) => s.enrollments);
  const setSyncing = useEnrollmentStore((s) => s.setSyncing);
  const updateStatus = useEnrollmentStore((s) => s.updateSyncStatus);
  const isSyncing = useEnrollmentStore((s) => s.isSyncing);
  const [result, setResult] = useState<{ synced: number; failed: number } | null>(null);

  const pending = enrollments.filter((e) => e.syncStatus === 'pending' || e.syncStatus === 'failed');

  const handleSync = async () => {
    setSyncing(true);
    setResult(null);
    let synced = 0, failed = 0;

    for (const enrollment of pending) {
      updateStatus(enrollment.id, 'syncing');
      try {
        // Simulate API call
        await new Promise((resolve) => setTimeout(resolve, 500));
        updateStatus(enrollment.id, 'synced');
        synced++;
      } catch {
        updateStatus(enrollment.id, 'failed');
        failed++;
      }
    }

    setSyncing(false);
    setResult({ synced, failed });
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Synchronisation</Text>
      <Text style={styles.subtitle}>{pending.length} inscription(s) en attente</Text>

      {isSyncing ? (
        <View style={styles.syncingBox}>
          <ActivityIndicator size="large" color="#065F46" />
          <Text style={styles.syncingText}>Synchronisation en cours...</Text>
        </View>
      ) : result ? (
        <View style={styles.resultBox}>
          <Text style={styles.resultTitle}>Synchronisation terminée</Text>
          <Text style={styles.resultDetail}>{result.synced} réussi(s), {result.failed} échoué(s)</Text>
        </View>
      ) : (
        <TouchableOpacity style={[styles.syncBtn, pending.length === 0 && styles.disabledBtn]}
          onPress={handleSync} disabled={pending.length === 0}>
          <Text style={styles.syncBtnText}>
            {pending.length === 0 ? 'Tout est synchronisé' : `Synchroniser ${pending.length} inscription(s)`}
          </Text>
        </TouchableOpacity>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f9fafb', padding: 20, justifyContent: 'center' },
  title: { fontSize: 24, fontWeight: 'bold', color: '#065F46', textAlign: 'center' },
  subtitle: { fontSize: 14, color: '#6b7280', textAlign: 'center', marginBottom: 30 },
  syncingBox: { alignItems: 'center', padding: 40 },
  syncingText: { marginTop: 16, fontSize: 16, color: '#065F46' },
  resultBox: { backgroundColor: '#ecfdf5', padding: 24, borderRadius: 12, alignItems: 'center' },
  resultTitle: { fontSize: 18, fontWeight: 'bold', color: '#065F46' },
  resultDetail: { fontSize: 14, color: '#6b7280', marginTop: 8 },
  syncBtn: { backgroundColor: '#D97706', padding: 16, borderRadius: 12, alignItems: 'center' },
  disabledBtn: { backgroundColor: '#d1d5db' },
  syncBtnText: { color: '#fff', fontSize: 16, fontWeight: '600' },
});
