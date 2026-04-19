import React, { useState } from 'react';
import { View, Text, TextInput, ScrollView, TouchableOpacity, StyleSheet } from 'react-native';
import { useEnrollmentStore } from '../stores/enrollment-store';
import type { EnrollmentData } from '../types';

export default function NewEnrollmentScreen({ navigation }: any) {
  const addEnrollment = useEnrollmentStore((s) => s.addEnrollment);
  const [form, setForm] = useState({
    firstName: '', lastName: '', maidenName: '', dateOfBirth: '', placeOfBirth: '',
    gender: 'M', nationalId: '', passportNumber: '', phone: '', email: '',
    countryOfResidence: '', cityOfResidence: '', addressAbroad: '',
    provinceOfOrigin: '', communeOfOrigin: '',
  });

  const update = (field: string, value: string) => setForm((f) => ({ ...f, [field]: value }));

  const handleSave = () => {
    const enrollment: EnrollmentData = {
      id: `enr-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
      ...form,
      gender: form.gender as 'M' | 'F',
      fingerprintsCapured: false,
      syncStatus: 'pending',
      createdAt: new Date().toISOString(),
    };
    addEnrollment(enrollment);
    navigation.navigate('Dashboard');
  };

  const fields: [string, string, string][] = [
    ['firstName', 'Prénom *', 'text'], ['lastName', 'Nom *', 'text'], ['maidenName', 'Nom de jeune fille', 'text'],
    ['dateOfBirth', 'Date de naissance * (JJ/MM/AAAA)', 'text'], ['placeOfBirth', 'Lieu de naissance *', 'text'],
    ['nationalId', 'N° CNIB', 'text'], ['passportNumber', 'N° Passeport', 'text'],
    ['phone', 'Téléphone', 'phone-pad'], ['email', 'Email', 'email-address'],
    ['countryOfResidence', 'Pays de résidence *', 'text'], ['cityOfResidence', 'Ville', 'text'],
    ['addressAbroad', 'Adresse complète', 'text'],
    ['provinceOfOrigin', "Province d'origine", 'text'], ['communeOfOrigin', "Commune d'origine", 'text'],
  ];

  return (
    <ScrollView style={styles.container}>
      <Text style={styles.title}>Nouvelle Inscription</Text>

      <View style={styles.genderRow}>
        {['M', 'F'].map((g) => (
          <TouchableOpacity key={g} style={[styles.genderBtn, form.gender === g && styles.genderActive]}
            onPress={() => update('gender', g)}>
            <Text style={[styles.genderText, form.gender === g && styles.genderTextActive]}>
              {g === 'M' ? 'Masculin' : 'Féminin'}
            </Text>
          </TouchableOpacity>
        ))}
      </View>

      {fields.map(([key, label]) => (
        <View key={key} style={styles.field}>
          <Text style={styles.label}>{label}</Text>
          <TextInput style={styles.input} value={(form as any)[key]}
            onChangeText={(v) => update(key, v)} placeholder={label} />
        </View>
      ))}

      <TouchableOpacity style={styles.saveBtn} onPress={handleSave}>
        <Text style={styles.saveBtnText}>Enregistrer l&apos;inscription</Text>
      </TouchableOpacity>
      <View style={{ height: 40 }} />
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f9fafb', padding: 20 },
  title: { fontSize: 22, fontWeight: 'bold', color: '#065F46', marginBottom: 20 },
  field: { marginBottom: 12 },
  label: { fontSize: 13, color: '#374151', marginBottom: 4, fontWeight: '500' },
  input: { backgroundColor: '#fff', borderWidth: 1, borderColor: '#d1d5db', borderRadius: 8, padding: 12, fontSize: 14 },
  genderRow: { flexDirection: 'row', marginBottom: 16, gap: 8 },
  genderBtn: { flex: 1, padding: 12, borderRadius: 8, borderWidth: 1, borderColor: '#d1d5db', alignItems: 'center', backgroundColor: '#fff' },
  genderActive: { backgroundColor: '#065F46', borderColor: '#065F46' },
  genderText: { color: '#374151', fontWeight: '500' },
  genderTextActive: { color: '#fff' },
  saveBtn: { backgroundColor: '#065F46', padding: 16, borderRadius: 12, alignItems: 'center', marginTop: 20 },
  saveBtnText: { color: '#fff', fontSize: 16, fontWeight: '600' },
});
