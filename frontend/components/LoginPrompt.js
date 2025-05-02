export default function LoginPrompt({ token }) {
  return (
    <div className="bg-yellow-100 border-l-4 border-yellow-500 text-yellow-800 p-4 mt-4 rounded">
      <p className="mb-2">You need to be logged in to accept this invite.</p>
      <div className="flex justify-center gap-4">
        <a
          className="underline text-blue-600 font-medium"
          href={`/login?redirect=/invite?token=${token}`}
        >
          Log In
        </a>
        <a
          className="underline text-blue-600 font-medium"
          href={`/register?redirect=/invite?token=${token}`}
        >
          Register
        </a>
      </div>
    </div>
  )
}
