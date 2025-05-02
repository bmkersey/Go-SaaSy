import Link from "next/link"
import Image from "next/image"

export default function HomePage() {
  return (
    <div className="min-h-screen bg-gray-50 text-black flex flex-col items-center justify-center px-4 py-16">
      <div className="max-w-2xl text-center">
      <Image
        src="/banner.png"
        alt="Go-SaaSy Banner"
        width={1200}
        height={400}
        className="w-full object-cover mb-8 rounded"
      />
        <h1 className="text-4xl font-bold mb-4">Go-SaaSy</h1>
        <p className="text-lg mb-6">A modern SaaS starter kit built with Go, Postgres, Stripe, and Next.js.</p>

        <div className="flex justify-center gap-4 mb-10">
        <Link
          href="/register"
          className="bg-blue-600 text-white px-6 py-2 rounded shadow hover:bg-blue-700"
        >
          Get Started
        </Link>
        <Link
          href="/login"
          className="bg-blue-600 text-white px-6 py-2 rounded shadow hover:bg-blue-700"
        >
          Log In
        </Link>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-left text-sm">
          <div>
            <h3 className="font-semibold mb-1">🔐 Auth</h3>
            <p>JWT-based login and register with bcrypt password hashing.</p>
          </div>
          <div>
            <h3 className="font-semibold mb-1">🏢 Orgs + Roles</h3>
            <p>Multi-tenant setup with org owners, members, and admins.</p>
          </div>
          <div>
            <h3 className="font-semibold mb-1">💳 Stripe Billing</h3>
            <p>Pro and Ultimate plans, subscriptions, and webhook handling.</p>
          </div>
          <div>
            <h3 className="font-semibold mb-1">📬 Email System</h3>
            <p>Templates for welcome, reset password, invites, and billing updates.</p>
          </div>
        </div>
      </div>

      <footer className="mt-16 text-xs text-gray-600">
        <p>
          Built by <a href="https://github.com/bmkersey" className="underline">bmkersey</a> — MIT Licensed
        </p>
      </footer>
    </div>
  )
}

