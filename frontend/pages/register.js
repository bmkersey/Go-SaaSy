import { useState } from 'react'
import { useRouter } from 'next/router'
import Image from 'next/image'

export default function Register() {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState(null)
  const router = useRouter()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError(null)

    try {
      const res = await fetch('http://localhost:8080/api/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, email, password }),
      })

      if (!res.ok) {
        const body = await res.json()
        throw new Error(body.error || 'Registration failed')
      }

      const { token } = await res.json()
      document.cookie = `token=${token}; path=/; max-age=86400; samesite=lax`
      router.push('/dashboard')
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <div className="min-h-screen flex flex-col text-black items-center justify-center bg-gray-50 p-4">
      <Image 
              src="/sassywhite.png" alt="SaaSy logo" width={200} height={200} className="animate-pulse"
            />
      <form
        onSubmit={handleSubmit}
        className="bg-white p-6 rounded-lg shadow-md w-full max-w-md"
      >
        <h1 className="text-xl font-semibold mb-4 text-center">Register for SaaSy</h1>

        {error && <div className="text-red-600 mb-3">{error}</div>}

        <input
          type="text"
          placeholder="Full Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          className="w-full p-2 border mb-3 rounded"
        />

        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="w-full p-2 border mb-3 rounded"
        />

        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className="w-full p-2 border mb-4 rounded"
        />

        <button
          type="submit"
          className="bg-blue-600 text-white w-full py-2 rounded hover:bg-blue-700"
        >
          Register
        </button>
      </form>
    </div>
  )
}
