import { useState } from "react";
import { useRouter } from 'next/router';
import Image from "next/image";


export default function login() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState(null)
  const router = useRouter()
  const redirect = router.query.redirect || "/dashboard"


  const handleSubmit = async (e) => {
    e.preventDefault()
    setError(null)
  

    try {
      const res = await fetch("http://localhost:8080/api/login", {
        method: "POST",
        headers: {
          'Content-Type': "application/json"
        },
        body: JSON.stringify({email, password})
      })

      if (!res.ok) {
        const body = await res.json()
        throw new Error(body.error || "Log in failed")
      }

      const { token } = await res.json()
      document.cookie = `token=${token}; path=/; max-age=86400; samesite=lax`
      router.push(redirect)
    } catch (err) {
      setError(err.message)
    }
  }


  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 p-4">
      <Image 
        src="/sassywhite.png" alt="SaaSy logo" width={200} height={200} className="animate-pulse"
      />
      <form
        onSubmit={handleSubmit}
        className="bg-white p-6 rounded-lg shadow-md w-full max-w-md"
      >
        <h1 className="text-xl text-black font-semibold mb-4">Log in to SaaSy</h1>

        {error && <div className="text-red-600 mb-3">{error}</div>}

        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="w-full p-2 border mb-3 rounded text-black"
        />

        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className="w-full p-2 border mb-4 rounded text-black"
        />

        <button
          type="submit"
          className="bg-blue-600 text-white w-full py-2 rounded hover:bg-blue-700"
        >
          Log In
        </button>
      </form>
    </div>
  )
}