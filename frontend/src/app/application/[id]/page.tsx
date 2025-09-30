import Application from "@/components/Application/application";

interface ApplicationPageProps {
  params: Promise<{ id: string }>;
}

export default async function ApplicationPage({
  params,
}: ApplicationPageProps) {
  const { id } = await params;
  return <Application id={id} />;
}
