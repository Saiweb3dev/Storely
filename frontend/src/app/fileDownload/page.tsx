import FileDownload from "@/components/FileDownload";
import RecentUploads from "@/components/RecentUpload";
export default function Page() {
  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <FileDownload/>
      <RecentUploads/>
    </div>
  );
}
