import Link from 'next/link';
import { Button } from "./components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "./components/ui/card";
import { Shield, Users, BarChart3, Search } from "lucide-react";

export default function Home() {
  return (
    <div className="space-y-12">
      {/* Hero Section */}
      <div className="text-center space-y-6 py-12">
        <h1 className="text-5xl font-bold tracking-tight">
          Welcome to <span className="text-primary">Tundragon Corporation</span>
        </h1>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          Your EVE Online Intelligence Hub. Advanced character management, analytics, and combat tracking for serious corporations.
        </p>
        <div className="flex gap-4 justify-center">
          <Button asChild size="lg">
            <Link href="/dashboard">
              <BarChart3 className="mr-2 size-4" />
              View Dashboard
            </Link>
          </Button>
          <Button asChild variant="outline" size="lg">
            <Link href="/admin">
              <Shield className="mr-2 size-4" />
              Admin Panel
            </Link>
          </Button>
        </div>
      </div>

      {/* Features Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Users className="size-5 text-primary" />
              Character Management
            </CardTitle>
            <CardDescription>
              Track and manage multiple EVE Online characters with detailed analytics
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild variant="ghost" className="w-full">
              <Link href="/dashboard">Explore Characters</Link>
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Search className="size-5 text-primary" />
              Character Search
            </CardTitle>
            <CardDescription>
              Search for characters by name using our advanced EVE ESI integration
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild variant="ghost" className="w-full">
              <Link href="/dashboard">Search Now</Link>
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart3 className="size-5 text-primary" />
              Advanced Analytics
            </CardTitle>
            <CardDescription>
              Deep insights into killmails, ISK tracking, and combat performance
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild variant="ghost" className="w-full">
              <Link href="/dashboard">View Analytics</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}