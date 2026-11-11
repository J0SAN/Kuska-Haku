import { z } from "zod";

export const feedbackSchema = z.object({
  rating: z
    .number()
    .int()
    .min(1, "Elige una calificación de 1 a 5 estrellas.")
    .max(5, "La calificación máxima es 5 estrellas."),
  comment: z
    .string()
    .max(500, "El comentario no puede superar los 500 caracteres."),
});

export type FeedbackFormValues = z.infer<typeof feedbackSchema>;
