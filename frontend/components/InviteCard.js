export default function LoginPrompt({ orgName, onAccept }) {
  return (
<div className="min-h-screen text-black bg-gray-100 p-6 flex items-center justify-center">
      <div className="bg-white p-6 rounded shadow-md max-w-md w-full text-center">
        <h2 className="text-xl font-bold mb-2">You're invited to join SaaSy!</h2>
        <p className="mb-4">
          Organization: <strong>{orgName}</strong>
        </p>
        <button
          onClick={onAccept}
          className="bg-blue-600 text-white font-semibold px-4 py-2 rounded"
        >
          Accept Invite
        </button>
      </div>
    </div>
  )
}