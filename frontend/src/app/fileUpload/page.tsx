import FileUpload from "@/components/FileUpload";
import RecentUploads from "@/components/RecentUpload";
import { UploadsProvider } from '@/contexts/UploadsContext';

export default function Home() {
  return (
    <UploadsProvider>
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <FileUpload />
      <RecentUploads />
    </div>
  </UploadsProvider>
  );
}
