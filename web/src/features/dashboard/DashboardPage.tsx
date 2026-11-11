import { useEffect, useState } from "react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import httpClient, { getErrorMessage } from "@/shared/api/httpClient";
import type { DashboardData } from "./dashboard.types";

const penFormatter = new Intl.NumberFormat("es-PE", {
  style: "currency",
  currency: "PEN",
});

const PERIODS = [7, 14, 30];

export default function DashboardPage() {
  const [days, setDays] = useState(7);
  const [data, setData] = useState<DashboardData | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isCancelled = false;

    httpClient
      .get<DashboardData>(`/admin/dashboard?days=${days}`)
      .then((res) => {
        if (!isCancelled) {
          setData(res.data);
          setError(null);
        }
      })
      .catch((err) => {
        if (!isCancelled) {
          setError(getErrorMessage(err, "Error cargando el dashboard"));
        }
      });

    return () => {
      isCancelled = true;
    };
  }, [days]);

  const series = (data?.series ?? []).map((d) => ({
    ...d,
    label: d.date.slice(5),
  }));

  return (
    <section>
      <div className="section-header">
        <h2>Dashboard</h2>
        <div className="period-picker">
          {PERIODS.map((p) => (
            <button
              key={p}
              type="button"
              className={p === days ? "tab active" : "tab"}
              onClick={() => setDays(p)}
            >
              {p} días
            </button>
          ))}
        </div>
      </div>

      {error && <p className="error">{error}</p>}
      {!data && !error && <p className="hint">Cargando estadísticas...</p>}

      {data && (
        <>
          <div className="kpi-grid">
            <KpiCard
              title="Pasajeros hoy"
              value={String(data.passengersToday)}
            />
            <KpiCard
              title="Ingresos hoy"
              value={penFormatter.format(data.incomeToday)}
            />
            <KpiCard title="Buses activos" value={String(data.activeBuses)} />
            <KpiCard title="Rutas activas" value={String(data.activeRoutes)} />
            <KpiCard
              title="Sugerencias nuevas"
              value={String(data.newSuggestions)}
            />
          </div>

          <div className="chart-grid">
            <div className="chart-card">
              <h3>Pasajeros por día</h3>
              <ResponsiveContainer width="100%" height={240}>
                <BarChart data={series}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#dfe7e0" />
                  <XAxis dataKey="label" fontSize={12} />
                  <YAxis allowDecimals={false} fontSize={12} />
                  <Tooltip />
                  <Bar
                    dataKey="passengers"
                    name="Pasajeros"
                    fill="#2e7d32"
                    radius={[4, 4, 0, 0]}
                  />
                </BarChart>
              </ResponsiveContainer>
            </div>

            <div className="chart-card">
              <h3>Ingresos por día (S/)</h3>
              <ResponsiveContainer width="100%" height={240}>
                <BarChart data={series}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#dfe7e0" />
                  <XAxis dataKey="label" fontSize={12} />
                  <YAxis fontSize={12} />
                  <Tooltip
                    formatter={(v) => penFormatter.format(Number(v ?? 0))}
                  />
                  <Bar
                    dataKey="amount"
                    name="Ingresos"
                    fill="#81c784"
                    radius={[4, 4, 0, 0]}
                  />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>

          <h3>Buses con más ingresos ({data.days} días)</h3>
          <table className="data-table">
            <thead>
              <tr>
                <th>Bus</th>
                <th>Pasajeros</th>
                <th>Ingresos</th>
              </tr>
            </thead>
            <tbody>
              {data.topBuses.length === 0 && (
                <tr>
                  <td colSpan={3} className="empty">
                    Todavía no hay registros de ingresos en este período.
                  </td>
                </tr>
              )}
              {data.topBuses.map((b) => (
                <tr key={b.busCode}>
                  <td>{b.busCode}</td>
                  <td>{b.passengers}</td>
                  <td>{penFormatter.format(b.amount)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
    </section>
  );
}

function KpiCard({ title, value }: { title: string; value: string }) {
  return (
    <div className="kpi-card">
      <span className="kpi-title">{title}</span>
      <span className="kpi-value">{value}</span>
    </div>
  );
}
