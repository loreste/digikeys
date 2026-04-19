import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import DashboardScreen from '../screens/DashboardScreen';
import NewEnrollmentScreen from '../screens/NewEnrollmentScreen';
import EnrollmentListScreen from '../screens/EnrollmentListScreen';
import SyncScreen from '../screens/SyncScreen';

const Stack = createNativeStackNavigator();

export default function Navigation() {
  return (
    <NavigationContainer>
      <Stack.Navigator initialRouteName="Dashboard"
        screenOptions={{ headerStyle: { backgroundColor: '#065F46' }, headerTintColor: '#fff', headerTitleStyle: { fontWeight: 'bold' } }}>
        <Stack.Screen name="Dashboard" component={DashboardScreen} options={{ title: 'DIGIKEYS' }} />
        <Stack.Screen name="NewEnrollment" component={NewEnrollmentScreen} options={{ title: 'Nouvelle Inscription' }} />
        <Stack.Screen name="EnrollmentList" component={EnrollmentListScreen} options={{ title: 'Inscriptions' }} />
        <Stack.Screen name="Sync" component={SyncScreen} options={{ title: 'Synchronisation' }} />
      </Stack.Navigator>
    </NavigationContainer>
  );
}
