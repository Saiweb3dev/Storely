// types/auth.ts
export interface SignInFormData {
  email: string;
  password: string;
}

export interface SignUpFormData {
  userID : string;
  username: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: {
    userID: string;
    username: string;
    email: string;
    storageUsed: number;
    storageLimit: number;
  };
  message?: string;
}