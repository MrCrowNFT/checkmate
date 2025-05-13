export interface safeCredential {
  id: number;
  user_id: string;
  platform: string;
  name: string;
  created_at: Date;
}

export interface platformCredentialInput {
  platform: string;
  name: string;
  apiKey: string;
}
