import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'

export default function AdminDashboard() {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState(null)
  const [error, setError] = useState(null)
  const router = useRouter()

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await fetch("http://localhost:8080/api/admin/overview", {
          method: "GET",
          credentials: "include"
        })

        if (!res.ok) {
          if (res.status === 403) {
            router.push("/dashboard")
          }
          throw new Error("Not authorized")
        }

        const json = await res.json()
        setData(json)
        setLoading(false)
      } catch (err) {
        setError(err.message)
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  if (loading) return <div className="p-4">Loading admin data...</div>
  if (error) return <div className="p-4 text-red-600">Error: {error}</div>
  if (!data) return <div className="p-4 text-gray-600">No data available.</div>

  const users = data.users || []
  const orgs = data.orgs || []

  return (
    <div className="min-h-screen text-black bg-gray-100 p-6">
      <h1 className="text-2xl font-bold mb-6 text-center">Admin Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Orgs */}
        <div className="bg-white p-4 shadow rounded">
          <h2 className="text-xl font-semibold mb-3">Organizations ({orgs.length})</h2>
          <ul className="space-y-2 text-sm">
            {orgs.map((org) => (
              <li key={org.id} className="border-b pb-2">
                <strong>{org.name}</strong> — Plan: {org.plan || "N/A"} | Paid: {org.is_paid ? "✅" : "❌"}
              </li>
            ))}
          </ul>
        </div>

        {/* Users */}
        <div className="bg-white p-4 shadow rounded">
          <h2 className="text-xl font-semibold mb-3">Users ({users.length})</h2>
          <ul className="space-y-2 text-sm">
            {users.map((u) => (
              <li key={u.id} className="border-b pb-2">
                {u.email} {u.is_admin && <span className="text-xs text-green-600">(admin)</span>}
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  )
}
