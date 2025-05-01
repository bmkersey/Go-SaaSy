import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'

export default function Dashboard() {
  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState(null)
  const router = useRouter()

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const res = await fetch("http://localhost:8080/api/me", {
          method: "GET",
          credentials: "include", // 🔑 Send cookies
        })

        if (!res.ok) {
          throw new Error("Unauthorized")
        }

        const data = await res.json()
        setUser(data)
        setLoading(false)
      } catch (err) {
        router.push("/login")
      }
    }

    fetchUser()
  }, [])

  if (loading) return <div className="p-4">Loading...</div>

  return (
    <div className="min-h-screen text-black flex flex-col items-center justify-start bg-gray-50 p-4">
      <button
        onClick={() => {
          document.cookie = 'token=; path=/; max-age=0'
          router.push('/login')
        }}
        className="mt-6 bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700 self-end"
      >
        Log Out
      </button>
      <h1 className="text-2xl font-bold mb-2">Welcome, {user.email}</h1>
      <p className="text-sm text-gray-600">User ID: <code>{user.id}</code></p>
    </div>
  )
}
