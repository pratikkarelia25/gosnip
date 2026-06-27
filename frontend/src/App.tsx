import { useState, useEffect } from "react"
import { Link2, Copy, Check, ArrowRight, Clock } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { shortenUrl, getRedirectUrl, type ShortenResponse } from "@/api"

export default function App() {
  const [longUrl, setLongUrl] = useState("")
  const [shortCode, setShortCode] = useState("")
  const [expiresIn, setExpiresIn] = useState("")
  const [result, setResult] = useState<ShortenResponse | null>(null)
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)
  const [copied, setCopied] = useState(false)
  const [notFound, setNotFound] = useState(false)

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    if (params.get("error") === "not_found") {
      setNotFound(true)
      window.history.replaceState({}, "", "/")
    }
  }, [])

  async function handleSubmit(e: React.SubmitEvent) {
    e.preventDefault()
    setError("")
    setResult(null)
    setLoading(true)

    try {
      const data = await shortenUrl({
        long_url: longUrl,
        short_code: shortCode || undefined,
        expires_in_seconds: expiresIn ? parseInt(expiresIn) : undefined,
      })
      setResult(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong")
    } finally {
      setLoading(false)
    }
  }

  async function handleCopy() {
    if (!result) return
    await navigator.clipboard.writeText(getRedirectUrl(result.short_code))
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const shortLink = result ? getRedirectUrl(result.short_code) : ""
  const isExisting = result?.message.includes("already exists")

  return (
    <main className="min-h-screen bg-background flex items-center justify-center p-4">
      <div className="w-full max-w-lg space-y-6">

        <div className="text-center space-y-2">
          <div className="inline-flex items-center justify-center w-12 h-12 rounded-2xl bg-primary text-primary-foreground mb-2">
            <Link2 className="w-6 h-6" />
          </div>
          <h1 className="text-3xl font-bold tracking-tight">GoSnip</h1>
          <p className="text-muted-foreground">Shorten your links, simply.</p>
        </div>

        {notFound && (
          <div className="rounded-lg border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive flex items-center justify-between">
            <span>That short link doesn't exist or has expired.</span>
            <button onClick={() => setNotFound(false)} className="ml-4 opacity-60 hover:opacity-100 text-base leading-none">&times;</button>
          </div>
        )}

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Shorten a URL</CardTitle>
            <CardDescription>Paste a long URL and get a short link instantly.</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="long_url">Long URL</Label>
                <Input
                  id="long_url"
                  type="url"
                  placeholder="https://example.com/very/long/url"
                  value={longUrl}
                  onChange={e => setLongUrl(e.target.value)}
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-2">
                  <Label htmlFor="short_code">
                    Custom code <span className="text-muted-foreground text-xs">(optional)</span>
                  </Label>
                  <Input
                    id="short_code"
                    placeholder="my-link"
                    value={shortCode}
                    onChange={e => setShortCode(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="expires_in">
                    Expires in <span className="text-muted-foreground text-xs">(seconds)</span>
                  </Label>
                  <Input
                    id="expires_in"
                    type="number"
                    placeholder="3600"
                    min={1}
                    value={expiresIn}
                    onChange={e => setExpiresIn(e.target.value)}
                  />
                </div>
              </div>

              {error && (
                <p className="text-sm text-destructive">{error}</p>
              )}

              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? "Shortening..." : (
                  <span className="flex items-center gap-2">
                    Shorten <ArrowRight className="w-4 h-4" />
                  </span>
                )}
              </Button>
            </form>
          </CardContent>
        </Card>

        {result && (
          <Card className="border-primary/20 bg-primary/5">
            <CardContent className="pt-6 space-y-3">
              <div className="flex items-center gap-2">
                <span className="text-sm font-medium">Your short link</span>
                {isExisting && (
                  <Badge variant="secondary" className="text-xs">already existed</Badge>
                )}
              </div>

              <div className="flex items-center gap-2">
                <code className="flex-1 text-sm bg-background border rounded-md px-3 py-2 truncate">
                  {shortLink}
                </code>
                <Button size="icon" variant="outline" onClick={handleCopy}>
                  {copied ? <Check className="w-4 h-4 text-green-500" /> : <Copy className="w-4 h-4" />}
                </Button>
                <a
                  href={shortLink}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center justify-center h-9 w-9 rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground transition-colors"
                >
                  <Link2 className="w-4 h-4" />
                </a>
              </div>

              {result.expires_at && (
                <p className="text-xs text-muted-foreground flex items-center gap-1">
                  <Clock className="w-3 h-3" />
                  Expires {new Date(result.expires_at).toLocaleString()}
                </p>
              )}
            </CardContent>
          </Card>
        )}
      </div>
    </main>
  )
}
