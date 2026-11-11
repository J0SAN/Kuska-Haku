import { Controller, useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useParams } from "react-router-dom";
import { CheckCircle2, QrCode } from "lucide-react";
import { getErrorMessage, getErrorStatus } from "@/shared/api/httpClient";
import { Alert, AlertDescription } from "@/shared/ui/alert";
import { Button } from "@/shared/ui/button";
import { Card, CardContent } from "@/shared/ui/card";
import { Label } from "@/shared/ui/label";
import { Skeleton } from "@/shared/ui/skeleton";
import { Textarea } from "@/shared/ui/textarea";
import { usePublicVehicle, useSubmitFeedback } from "./hooks/useVehicleFeedback";
import { getDeviceId } from "./lib/deviceId";
import {
  feedbackSchema,
  type FeedbackFormValues,
} from "./vehicleFeedback.schema";
import { StarRatingInput } from "./components/StarRatingInput";

const MAX_COMMENT_LENGTH = 500;

export default function PublicFeedbackPage() {
  const { token } = useParams<{ token: string }>();
  const vehicle = usePublicVehicle(token);
  const submit = useSubmitFeedback(token);

  const {
    control,
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FeedbackFormValues>({
    resolver: zodResolver(feedbackSchema),
    defaultValues: { rating: 0, comment: "" },
  });

  const comment = useWatch({ control, name: "comment" }) ?? "";

  const onSubmit = (values: FeedbackFormValues) => {
    if (!token) return;
    submit.mutate({
      rating: values.rating,
      comment: values.comment,
      deviceId: getDeviceId(),
    });
  };

  const status = getErrorStatus(submit.error);
  const alreadySubmitted = status === 409;
  const networkError = submit.isError && !alreadySubmitted;

  return (
    <div className="min-h-screen bg-background px-4 py-6">
      {/* React 19 eleva estas etiquetas al <head>. */}
      <title>Califica tu vehículo · Kuska Haku</title>
      <meta name="robots" content="noindex" />

      <main className="mx-auto w-full max-w-md space-y-4">
        <header className="space-y-1">
          <p className="text-sm font-semibold text-secondary-foreground">
            Kuska Haku
          </p>
          <h2 className="text-2xl">Califica este vehículo</h2>
        </header>

        {vehicle.isLoading && (
          <Card>
            <CardContent className="space-y-3 pt-2">
              <Skeleton className="h-5 w-3/4" />
              <Skeleton className="h-5 w-1/2" />
              <Skeleton className="h-32 w-full" />
            </CardContent>
          </Card>
        )}

        {vehicle.isError && (
          <Card>
            <CardContent className="flex flex-col items-center gap-3 py-8 text-center">
              <QrCode aria-hidden className="size-10 text-muted-foreground" />
              <p className="text-base font-medium">
                {getErrorStatus(vehicle.error) === 404
                  ? "Este código QR ya no está disponible."
                  : "No pudimos cargar el vehículo."}
              </p>
              {getErrorStatus(vehicle.error) !== 404 && (
                <Button type="button" onClick={() => void vehicle.refetch()}>
                  Reintentar
                </Button>
              )}
            </CardContent>
          </Card>
        )}

        {vehicle.data && submit.isSuccess && (
          <Card>
            <CardContent className="flex flex-col items-center gap-3 py-10 text-center">
              <CheckCircle2 aria-hidden className="size-12 text-primary" />
              <p className="text-lg font-semibold">¡Gracias por tu opinión!</p>
              <p className="text-sm text-muted-foreground">
                Tu calificación de {vehicle.data.busCode} quedó registrada.
              </p>
            </CardContent>
          </Card>
        )}

        {vehicle.data && !submit.isSuccess && alreadySubmitted && (
          <Card>
            <CardContent className="flex flex-col items-center gap-3 py-10 text-center">
              <CheckCircle2 aria-hidden className="size-12 text-primary" />
              <p className="text-lg font-semibold">
                Ya calificaste este vehículo hoy. ¡Gracias!
              </p>
            </CardContent>
          </Card>
        )}

        {vehicle.data && !submit.isSuccess && !alreadySubmitted && (
          <Card>
            <CardContent className="space-y-5 pt-2">
              <div className="rounded-md bg-muted px-3 py-2 text-sm">
                Bus <strong>{vehicle.data.busCode}</strong> · Ruta{" "}
                <strong>{vehicle.data.routeName}</strong>
              </div>

              <form
                className="space-y-5"
                onSubmit={handleSubmit(onSubmit)}
                noValidate
              >
                <Controller
                  control={control}
                  name="rating"
                  render={({ field }) => (
                    <div className="space-y-1">
                      <StarRatingInput
                        value={field.value}
                        onChange={field.onChange}
                      />
                      {errors.rating && (
                        <p className="text-sm text-destructive" role="alert">
                          {errors.rating.message}
                        </p>
                      )}
                    </div>
                  )}
                />

                <div className="space-y-2">
                  <Label htmlFor="comment">Comentario (opcional)</Label>
                  <Textarea
                    id="comment"
                    rows={4}
                    maxLength={MAX_COMMENT_LENGTH}
                    placeholder="Cuéntanos cómo fue tu viaje"
                    aria-describedby="comment-counter"
                    {...register("comment")}
                  />
                  <div className="flex items-center justify-between text-xs text-muted-foreground">
                    <span id="comment-counter">
                      {comment.length}/{MAX_COMMENT_LENGTH}
                    </span>
                    {errors.comment && (
                      <span className="text-destructive" role="alert">
                        {errors.comment.message}
                      </span>
                    )}
                  </div>
                </div>

                {networkError && (
                  <Alert variant="destructive">
                    <AlertDescription>
                      {getErrorMessage(
                        submit.error,
                        "No pudimos guardar tu calificación. Revisa tu conexión.",
                      )}
                    </AlertDescription>
                  </Alert>
                )}

                <Button
                  type="submit"
                  size="lg"
                  className="w-full"
                  disabled={submit.isPending}
                >
                  {submit.isPending
                    ? "Enviando…"
                    : networkError
                      ? "Reintentar"
                      : "Enviar"}
                </Button>

                <p className="text-xs text-muted-foreground">
                  Tu opinión es anónima. Evita escribir datos personales.
                </p>
              </form>
            </CardContent>
          </Card>
        )}
      </main>
    </div>
  );
}
