import Home from "@/components/Home/home";

interface HomePageProps {
  params: Promise<{ page: string }>;
}

export default async function HomePage({ params }: HomePageProps) {
  const { page } = await params;
  return <Home page={page} />;
}
