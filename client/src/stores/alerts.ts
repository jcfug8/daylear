import { defineStore } from 'pinia';

export type AlertSeverity = 'info' | 'success' | 'warning' | 'error';

export interface Alert {
  id: number;
  message: string;
  severity: AlertSeverity;
  timeoutId?: ReturnType<typeof setTimeout>;
}

export const useAlertStore = defineStore('alerts', {
  state: () => ({
    alerts: [] as Alert[],
    nextId: 1,
  }),
  actions: {
    addAlert(message: string, severity: AlertSeverity = 'info', duration = 5000) {
      const id = this.nextId++;
      const alert: Alert = { id, message, severity };
      // Auto-remove after duration
      alert.timeoutId = setTimeout(() => {
        this.removeAlert(id);
      }, duration);
      this.alerts.push(alert);
      return id;
    },
    removeAlert(id: number) {
      const idx = this.alerts.findIndex(a => a.id === id);
      if (idx !== -1) {
        const [alert] = this.alerts.splice(idx, 1);
        if (alert.timeoutId) clearTimeout(alert.timeoutId);
      }
    },
    clearAll() {
      this.alerts.forEach(alert => {
        if (alert.timeoutId) clearTimeout(alert.timeoutId);
      });
      this.alerts = [];
    },
  },
}); 