import Link from 'next/link';
import { Button } from "./components/ui/button";

export default function Home() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-8">
      <h1 className="text-4xl font-bold mb-8">EVE Ran</h1>
      <div className="flex gap-4">
        <Button asChild>
          <Link href="/dashboard">Dashboard</Link>
        </Button>
        <Button asChild variant="secondary">
          <Link href="/characters">Character Management</Link>
        </Button>
        <Button asChild variant="outline">
          <Link href="/charts">Performance Charts</Link>
        </Button>
      </div>
    </div>
  );
}