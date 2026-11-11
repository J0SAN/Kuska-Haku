import { useSearchParams } from "react-router-dom";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/shared/ui/tabs";
import { FeedbackTab } from "./components/FeedbackTab";
import { QrCodesTab } from "./components/QrCodesTab";

type TabValue = "qr" | "comments";

export default function VehicleFeedbackPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const tab: TabValue =
    searchParams.get("tab") === "comments" ? "comments" : "qr";

  const handleTabChange = (value: string) => {
    setSearchParams(value === "comments" ? { tab: "comments" } : {}, {
      replace: true,
    });
  };

  return (
    <section className="space-y-4">
      <div>
        <h2>Calificaciones QR</h2>
        <p className="text-sm text-muted-foreground">
          Genera el QR de cada vehículo y revisa lo que opinan los pasajeros.
        </p>
      </div>

      <Tabs value={tab} onValueChange={handleTabChange}>
        <TabsList>
          <TabsTrigger value="qr">Códigos QR</TabsTrigger>
          <TabsTrigger value="comments">Comentarios</TabsTrigger>
        </TabsList>
        <TabsContent value="qr">
          <QrCodesTab />
        </TabsContent>
        <TabsContent value="comments">
          <FeedbackTab />
        </TabsContent>
      </Tabs>
    </section>
  );
}
