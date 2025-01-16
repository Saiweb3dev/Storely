import FileUpload from "@/components/File/FileUpload";
import RecentUploads from "@/components/RecentUpload";

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <FileUpload />
      <RecentUploads />
    </div>
  );
}
