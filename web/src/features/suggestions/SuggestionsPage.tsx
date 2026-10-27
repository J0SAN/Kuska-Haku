import { useEffect, useState } from "react";
import apiClient, {
  getErrorMessage,
  type SuggestionData,
} from "../../api/client";

export default function SuggestionsPage() {
  const [list, setList] = useState<SuggestionData[]>([]);
  const [onlyNew, setOnlyNew] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadSuggestions = async () => {
      try {
        const query = onlyNew ? "?only=new" : "";
        const res = await apiClient.get<SuggestionData[]>(
          `/admin/suggestions${query}`,
        );
        if (isMounted) {
          setList(res.data || []);
          setError(null);
        }
      } catch (err) {
        if (isMounted) {
          setError(getErrorMessage(err, "Error cargando las sugerencias"));
        }
      }
    };

    void loadSuggestions();

    return () => {
      isMounted = false;
    };
  }, [onlyNew]);

  const handleMarkRead = async (id: string) => {
    try {
      await apiClient.put(`/admin/suggestions/${id}/read`);
      const query = onlyNew ? "?only=new" : "";
      const res = await apiClient.get<SuggestionData[]>(
        `/admin/suggestions${query}`,
      );
      setList(res.data || []);
      setError(null);
    } catch (err) {
      setError(getErrorMessage(err, "Error actualizando la sugerencia"));
    }
  };

  const formatDate = (iso: string) =>
    new Date(iso).toLocaleString("es-PE", {
      dateStyle: "medium",
      timeStyle: "short",
    });

  return (
    <section>
      <div className="section-header">
        <h2>Buzón de sugerencias</h2>
        <label className="check">
          <input
            type="checkbox"
            checked={onlyNew}
            onChange={(e) => setOnlyNew(e.target.checked)}
          />
          Mostrar solo nuevas
        </label>
      </div>

      {error && <p className="error">{error}</p>}

      <table className="data-table">
        <thead>
          <tr>
            <th>Fecha</th>
            <th>Usuario</th>
            <th>Mensaje</th>
            <th>Estado</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {list.length === 0 && (
            <tr>
              <td colSpan={5} className="empty">
                No hay sugerencias para mostrar.
              </td>
            </tr>
          )}
          {list.map((s) => (
            <tr key={s.id}>
              <td>{formatDate(s.createdAt)}</td>
              <td>{s.userName}</td>
              <td className="message-cell">{s.message}</td>
              <td>
                <span
                  className={s.status === "new" ? "badge badge-new" : "badge"}
                >
                  {s.status === "new" ? "Nueva" : "Leída"}
                </span>
              </td>
              <td>
                {s.status === "new" && (
                  <button
                    type="button"
                    className="secondary"
                    onClick={() => handleMarkRead(s.id)}
                  >
                    Marcar como leída
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}
