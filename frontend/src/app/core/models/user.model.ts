export type Role = 'admin' | 'patient' | 'doctor' | 'corporate_hr';

export interface User {
  id: string;
  email: string;
  name: string;
  phone: string;
  role: Role;
  is_verified: boolean;
}

export interface AuthResponse {
  success: boolean;
  data: {
    user: User;
    token: string;
  };
}
