export interface UserData {
  id: string;
  email: string;
  fullName?: string;
  full_name?: string;
  role: "passenger" | "driver" | "admin";
  createdAt?: string;
}
