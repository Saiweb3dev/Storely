import FileDownload from "@/components/File/FileDownload";
import RecentUploads from "@/components/RecentUpload";
export default function Page() {
  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <FileDownload/>
      <RecentUploads/>
    </div>
  );
}
