// Server-side base URL for reaching the backend from within the container/network.
const apiBase = process.env.API_BASE_URL ?? "http://localhost:12010";

async function getBackendStatus(): Promise<string> {
  try {
    const res = await fetch(`${apiBase}/uptime`, { cache: "no-store" });
    return res.ok ? await res.text() : `error: ${res.status}`;
  } catch {
    return "unreachable";
  }
}

export default async function Home() {
  const status = await getBackendStatus();

  return (
    <main className="mx-auto flex max-w-2xl flex-col gap-6 px-6 py-16">
      <header className="flex flex-col gap-2">
        <h1 className="text-3xl font-bold tracking-tight">go-github-tracker</h1>
        <p className="text-gray-600">
          Next.js + React + TypeScript + Tailwind frontend. Fill in the UI here.
        </p>
      </header>

      <section className="rounded-lg border border-gray-200 bg-white p-5 shadow-sm">
        <div className="text-sm font-medium text-gray-500">Backend API</div>
        <code className="text-sm text-gray-800">{apiBase}</code>
        <div className="mt-3 flex items-center gap-2 text-sm">
          <span className="text-gray-500">/uptime:</span>
          <span
            className={
              status === "ok"
                ? "rounded bg-green-100 px-2 py-0.5 font-mono text-green-700"
                : "rounded bg-red-100 px-2 py-0.5 font-mono text-red-700"
            }
          >
            {status}
          </span>
        </div>
      </section>
    </main>
  );
}
