import { useRef } from "react";
import { createPortal } from "react-dom";
import { QRCodeCanvas, QRCodeSVG } from "qrcode.react";
import { Download, Printer } from "lucide-react";
import { Alert, AlertDescription } from "@/shared/ui/alert";
import { Button } from "@/shared/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/dialog";
import { buildQrUrl, canPrintQr, getQrUrlWarning } from "../lib/qrUrl";
import type { VehicleQRSummary } from "../vehicleFeedback.types";

interface Props {
  item: VehicleQRSummary | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  link.click();
  URL.revokeObjectURL(url);
}

export function QrPreviewDialog({ item, open, onOpenChange }: Props) {
  const svgRef = useRef<HTMLDivElement>(null);
  const url = item?.token ? buildQrUrl(item.token) : null;
  const printable = canPrintQr();
  const warning = getQrUrlWarning();

  if (!item) return null;

  const canvasId = `qr-png-${item.fleetId}`;

  const handleDownloadSvg = () => {
    const svg = svgRef.current?.querySelector("svg");
    if (!svg) return;
    const serialized = new XMLSerializer().serializeToString(svg);
    downloadBlob(
      new Blob([`<?xml version="1.0" encoding="UTF-8"?>\n${serialized}`], {
        type: "image/svg+xml;charset=utf-8",
      }),
      `qr-${item.busCode}.svg`,
    );
  };

  const handleDownloadPng = () => {
    const canvas = document.getElementById(canvasId) as HTMLCanvasElement | null;
    if (!canvas) return;
    const link = document.createElement("a");
    link.href = canvas.toDataURL("image/png");
    link.download = `qr-${item.busCode}.png`;
    link.click();
  };

  const handlePrint = () => {
    window.print();
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>QR de {item.busCode}</DialogTitle>
          <DialogDescription>
            Ruta {item.routeName}. El pasajero lo escanea y califica solo este
            vehículo.
          </DialogDescription>
        </DialogHeader>

        {url ? (
          <div className="flex flex-col items-center gap-4">
            <div ref={svgRef} className="rounded-lg bg-card p-4">
              <QRCodeSVG
                value={url}
                size={220}
                level="Q"
                marginSize={4}
                bgColor="#ffffff"
                fgColor="#1b5e20"
                title={`QR del vehículo ${item.busCode}`}
              />
            </div>
            <p className="w-full break-all rounded-md bg-muted px-3 py-2 text-center text-xs text-muted-foreground">
              {url}
            </p>
            {/* Lienzo oculto a 1024 px solo para exportar el PNG. */}
            <div className="hidden" aria-hidden>
              <QRCodeCanvas
                id={canvasId}
                value={url}
                size={1024}
                level="Q"
                marginSize={4}
                bgColor="#ffffff"
                fgColor="#1b5e20"
              />
            </div>
          </div>
        ) : (
          <Alert variant="destructive">
            <AlertDescription>
              No hay una URL pública configurada para este QR.
            </AlertDescription>
          </Alert>
        )}

        {warning && (
          <Alert>
            <AlertDescription>{warning}</AlertDescription>
          </Alert>
        )}

        <p className="text-xs text-muted-foreground">
          Imprime el sticker de al menos 5 × 5 cm. El código del bus aparece en
          la hoja para que el pasajero compruebe que coincide.
        </p>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={handleDownloadSvg}
            disabled={!url || !printable}
          >
            <Download aria-hidden /> SVG
          </Button>
          <Button
            type="button"
            variant="outline"
            onClick={handleDownloadPng}
            disabled={!url || !printable}
          >
            <Download aria-hidden /> PNG
          </Button>
          <Button
            type="button"
            onClick={handlePrint}
            disabled={!url || !printable}
          >
            <Printer aria-hidden /> Imprimir
          </Button>
        </DialogFooter>

        {/* Hoja que se envía a la impresora (el resto se oculta con @media print).
            Va en un portal para que ningún contenedor del diálogo la recorte. */}
        {createPortal(
          <div className="print-sheet hidden print:block">
            <p className="text-center text-lg font-semibold">
              Escanea y califica este vehículo
            </p>
            {url && (
              <div className="mt-4 flex justify-center">
                <QRCodeSVG
                  value={url}
                  size={260}
                  level="Q"
                  marginSize={4}
                  bgColor="#ffffff"
                  fgColor="#000000"
                />
              </div>
            )}
            <p className="mt-4 text-center text-2xl font-bold">{item.busCode}</p>
            <p className="text-center text-sm">Kuska Haku</p>
          </div>,
          document.body,
        )}
      </DialogContent>
    </Dialog>
  );
}
