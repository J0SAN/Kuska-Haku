interface Props {
  title: string;
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
}

export function ConfirmDialog({ title, message, onConfirm, onCancel }: Props) {
  return (
    <div className="dialog-backdrop" role="dialog" aria-modal="true">
      <div className="dialog">
        <h3>{title}</h3>
        <p>{message}</p>
        <div className="dialog-actions">
          <button type="button" className="secondary" onClick={onCancel}>
            Cancelar
          </button>
          <button type="button" className="danger" onClick={onConfirm}>
            Eliminar
          </button>
        </div>
      </div>
    </div>
  );
}
