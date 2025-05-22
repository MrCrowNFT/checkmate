export interface safeCredential {
  id: number;
  user_id: string;
  platform: string;
  created_at: Date;
}

export interface platformCredentialInput {
  platform: string;
  apiKey: string;
}
