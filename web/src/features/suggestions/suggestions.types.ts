export interface SuggestionData {
  id: string;
  userId: string;
  userName: string;
  message: string;
  status: "new" | "read";
  createdAt: string;
}
