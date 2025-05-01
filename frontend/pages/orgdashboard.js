import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'

export default function OrgDashboard() {
  const [loading, setLoading] = useState(true)
  const [org, setOrg] = useState(null)
  const [members, setMembers] = useState([])
  const [error, setError] = useState(null)
  const router = useRouter()

  useEffect(() => {
    const fetchDashboard = async () => {
      try {
        const res = await fetch("http://localhost:8080/api/org/dashboard", {
          method: "GET",
          credentials: "include",
        })

        if (!res.ok) {
          throw new Error("Not authorized or not the org owner.")
        }

        const data = await res.json()
        setOrg(data.org)
        setMembers(data.members)
        setLoading(false)
      } catch (err) {
        setError(err.message)
        setLoading(false)
      }
    }

    fetchDashboard()
  }, [])

  if (loading) return <div className="p-4">Loading...</div>
  if (error) return <div className="p-4 text-red-600">Error: {error}</div>

  return (
    <div className="min-h-screen text-black bg-gray-100 p-6">
      <div className="max-w-3xl mx-auto space-y-6">

        <div className="bg-white shadow-md rounded p-6">
          <h1 className="text-2xl font-bold mb-4">Organization Dashboard</h1>
          <p><span className="font-medium">Name:</span> {org.name}</p>
          <p><span className="font-medium">Plan:</span> {org.plan || "N/A"}</p>
          <p><span className="font-medium">Billing Email:</span> {org.billing_email || "N/A"}</p>
          <p><span className="font-medium">Is Paid:</span> {org.is_paid ? "Yes" : "No"}</p>
        </div>

        <div className="bg-white shadow-md rounded p-6">
          <h2 className="text-xl font-semibold mb-4">Team Members ({members.length})</h2>
          <ul className="space-y-2 text-sm">
            {members.map((m) => (
              <li key={m.id} className="border-b pb-2">{m.email}</li>
            ))}
          </ul>
        </div>

      </div>
    </div>
  )
}
