import { useEffect, useState } from "react"
import { useRouter } from "next/router"
import LoginPrompt from "@/components/LoginPrompt"
import InviteCard from "@/components/InviteCard"

export default function InvitePage() {
  const router = useRouter()
  const token = router.query.token
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [inviteInfo, setInviteInfo] = useState(null)
  const [accepted, setAccepted] = useState(false)
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  useEffect(() => {
    fetch("http://localhost:8080/api/me", { credentials: "include" })
      .then((res) => {
        if (res.ok) setIsLoggedIn(true)
      })
      .catch(() => {})
  }, [])

  useEffect(() => {
    if (!token) return

    const validate = async () => {
      try {
        const res = await fetch(`http://localhost:8080/api/invite/validate?token=${token}`, {
          credentials: "include",
        })

        if (!res.ok) throw new Error("Invalid or expired invite")

        const data = await res.json()
        setInviteInfo(data)
        setLoading(false)
      } catch (err) {
        setError(err.message)
        setLoading(false)
      }
    }

    validate()
  }, [token])

  const handleAccept = async () => {
    try {
      const res = await fetch("http://localhost:8080/api/invite/accept", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ token }),
      })

      if (!res.ok) throw new Error("Failed to accept invite")

      setAccepted(true)
      setTimeout(() => router.push("/dashboard"), 2000)
    } catch (err) {
      setError(err.message)
    }
  }

  if (loading) return <div className="p-4">Validating invite...</div>
  if (error === "unauthorized") {
    return (
      <div className="p-4">
        <LoginPrompt token={token} />
      </div>
    )
  }
  if (accepted) return <div className="p-4 text-green-600">✅ Invite accepted! Redirecting...</div>

  return isLoggedIn ? (
    <InviteCard orgName={inviteInfo.org_name} onAccept={handleAccept} />
  ) : (
    <LoginPrompt token={token} />
  )
}
