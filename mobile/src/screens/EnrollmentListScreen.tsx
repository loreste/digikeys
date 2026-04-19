import React from 'react';
import { View, Text, FlatList, StyleSheet } from 'react-native';
import { useEnrollmentStore } from '../stores/enrollment-store';

const statusColors: Record<string, string> = {
  pending: '#D97706', synced: '#059669', failed: '#DC2626', syncing: '#6366F1',
};

export default function EnrollmentListScreen() {
  const enrollments = useEnrollmentStore((s) => s.enrollments);

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Inscriptions ({enrollments.length})</Text>
      <FlatList
        data={enrollments}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <View style={styles.card}>
            <View style={styles.cardHeader}>
              <Text style={styles.name}>{item.firstName} {item.lastName}</Text>
              <View style={[styles.badge, { backgroundColor: statusColors[item.syncStatus] || '#9ca3af' }]}>
                <Text style={styles.badgeText}>{item.syncStatus}</Text>
              </View>
            </View>
            <Text style={styles.detail}>{item.countryOfResidence} - {item.dateOfBirth}</Text>
            <Text style={styles.date}>{new Date(item.createdAt).toLocaleDateString('fr-FR')}</Text>
          </View>
        )}
        ListEmptyComponent={<Text style={styles.empty}>Aucune inscription enregistrée</Text>}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f9fafb', padding: 20 },
  title: { fontSize: 22, fontWeight: 'bold', color: '#065F46', marginBottom: 16 },
  card: { backgroundColor: '#fff', padding: 16, borderRadius: 12, marginBottom: 8, borderLeftWidth: 3, borderLeftColor: '#065F46' },
  cardHeader: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 4 },
  name: { fontSize: 16, fontWeight: '600', color: '#1f2937' },
  badge: { paddingHorizontal: 8, paddingVertical: 2, borderRadius: 10 },
  badgeText: { color: '#fff', fontSize: 10, fontWeight: '600' },
  detail: { fontSize: 13, color: '#6b7280' },
  date: { fontSize: 11, color: '#9ca3af', marginTop: 4 },
  empty: { textAlign: 'center', color: '#9ca3af', marginTop: 40 },
});
