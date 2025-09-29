import Report from "@/components/Report/report";

interface ReportPageProps {
  params: Promise<{ id: string; page: string }>;
}

export default async function ReportPage({ params }: ReportPageProps) {
  const { id, page } = await params;
  return <Report id={id} page={page} />;
}
