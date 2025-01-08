export interface SignInFormData {
  email: string
  password: string
}

export interface UserData {
  name: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  message?: string;
}

export interface UserPreferences {
  theme: 'light' | 'dark';
  notifications: boolean;
}

export interface ProfileMetadata {
  profilePictureURL?: string;
  preferences: UserPreferences;
}

export interface SignUpFormData {
  username: string;
  email: string;
  password: string;
  country: string;
  profilePicture?: File;
  preferences: {
    theme: 'light' | 'dark';
    notifications: boolean;
  };
}