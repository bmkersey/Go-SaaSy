import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'

export default function Dashboard() {
  const [loading, setLoading] = useState(true)
  const [token, setToken] = useState(null)
  const router = useRouter()

  // Helper to read token from cookies
  function getCookie(name) {
    const value = `; ${document.cookie}`
    const parts = value.split(`; ${name}=`)
    if (parts.length === 2) return parts.pop().split(';').shift()
  }

  useEffect(() => {
    const t = getCookie('token')
    if (!t) {
      router.push('/login')
    } else {
      setToken(t)
      setLoading(false)
    }
  }, [])

  if (loading) return <div className="p-4">Loading...</div>

  return (
    <div className="min-h-screen flex flex-col items-center justify-start bg-gray-50">
      <button
  onClick={() => {
    // Clear cookie by setting it to expire
    document.cookie = 'token=; path=/; max-age=0'
    router.push('/login')
  }}
  className="mt-6 bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700 self-end mr-10"
>
  Log Out
</button>
      <h1 className="text-2xl font-bold text-black mb-2">Welcome to the Dashboard</h1>
      <h2 className='text-black font-bold'>Auth is handled by cookies and tokens</h2>
      <p className="text-sm text-gray-600 w-full overflow-x-auto px-10 whitespace-nowrap">
  Your token: <code>{token}</code>
</p>
    </div>
  )
}
