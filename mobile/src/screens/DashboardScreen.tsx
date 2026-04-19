import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { useEnrollmentStore } from '../stores/enrollment-store';

export default function DashboardScreen({ navigation }: any) {
  const enrollments = useEnrollmentStore((s) => s.enrollments);
  const pending = enrollments.filter((e) => e.syncStatus === 'pending').length;
  const synced = enrollments.filter((e) => e.syncStatus === 'synced').length;
  const failed = enrollments.filter((e) => e.syncStatus === 'failed').length;

  return (
    <View style={styles.container}>
      <Text style={styles.title}>DIGIKEYS Enrôlement</Text>
      <Text style={styles.subtitle}>Agent d&apos;enrôlement mobile</Text>

      <View style={styles.statsRow}>
        <View style={[styles.statCard, { borderLeftColor: '#D97706' }]}>
          <Text style={styles.statValue}>{pending}</Text>
          <Text style={styles.statLabel}>En attente</Text>
        </View>
        <View style={[styles.statCard, { borderLeftColor: '#059669' }]}>
          <Text style={styles.statValue}>{synced}</Text>
          <Text style={styles.statLabel}>Synchronisés</Text>
        </View>
        <View style={[styles.statCard, { borderLeftColor: '#DC2626' }]}>
          <Text style={styles.statValue}>{failed}</Text>
          <Text style={styles.statLabel}>Échoués</Text>
        </View>
      </View>

      <TouchableOpacity style={styles.primaryBtn} onPress={() => navigation.navigate('NewEnrollment')}>
        <Text style={styles.primaryBtnText}>Nouvelle Inscription</Text>
      </TouchableOpacity>

      <TouchableOpacity style={styles.secondaryBtn} onPress={() => navigation.navigate('EnrollmentList')}>
        <Text style={styles.secondaryBtnText}>Voir les inscriptions ({enrollments.length})</Text>
      </TouchableOpacity>

      <TouchableOpacity style={styles.syncBtn} onPress={() => navigation.navigate('Sync')}>
        <Text style={styles.syncBtnText}>Synchroniser ({pending} en attente)</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f9fafb', padding: 20 },
  title: { fontSize: 28, fontWeight: 'bold', color: '#065F46', textAlign: 'center', marginTop: 40 },
  subtitle: { fontSize: 14, color: '#6b7280', textAlign: 'center', marginBottom: 30 },
  statsRow: { flexDirection: 'row', justifyContent: 'space-between', marginBottom: 30 },
  statCard: { flex: 1, backgroundColor: '#fff', padding: 16, marginHorizontal: 4, borderRadius: 12, borderLeftWidth: 4, alignItems: 'center' },
  statValue: { fontSize: 24, fontWeight: 'bold', color: '#1f2937' },
  statLabel: { fontSize: 11, color: '#6b7280', marginTop: 4 },
  primaryBtn: { backgroundColor: '#065F46', padding: 16, borderRadius: 12, alignItems: 'center', marginBottom: 12 },
  primaryBtnText: { color: '#fff', fontSize: 16, fontWeight: '600' },
  secondaryBtn: { backgroundColor: '#fff', padding: 16, borderRadius: 12, alignItems: 'center', marginBottom: 12, borderWidth: 1, borderColor: '#d1d5db' },
  secondaryBtnText: { color: '#374151', fontSize: 14 },
  syncBtn: { backgroundColor: '#D97706', padding: 16, borderRadius: 12, alignItems: 'center' },
  syncBtnText: { color: '#fff', fontSize: 14, fontWeight: '600' },
});
