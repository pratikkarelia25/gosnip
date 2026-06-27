const BASE_URL = import.meta.env.VITE_API_URL as string

export interface ShortenRequest {
  long_url: string
  short_code?: string
  expires_in_seconds?: number
}

export interface ShortenResponse {
  message: string
  short_code: string
  expires_at?: string
}

export async function shortenUrl(body: ShortenRequest): Promise<ShortenResponse> {
  const res = await fetch(`${BASE_URL}/shorten`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  })

  const data = await res.json()

  if (!res.ok) {
    throw new Error(data.error ?? "Something went wrong")
  }

  return data as ShortenResponse
}

export function getRedirectUrl(code: string): string {
  return `${BASE_URL}/${code}`
}
