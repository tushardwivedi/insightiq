import { Button } from "@/components/ui/button";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface PaginationProps {
  page: number;
  onPageChange: (page: number) => void;
  hasMore: boolean;
  totalShowing?: number;
  label?: string;
}

export function Pagination({ page, onPageChange, hasMore, totalShowing, label = "items" }: PaginationProps) {
  return (
    <div className="flex items-center justify-between mt-4 pt-4 border-t">
      {totalShowing !== undefined && (
        <p className="text-sm text-muted-foreground">
          Showing {totalShowing} {label}
        </p>
      )}
      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange(Math.max(0, page - 1))}
          disabled={page === 0}
        >
          <ChevronLeft className="w-4 h-4 mr-1" />
          Previous
        </Button>
        <span className="text-sm text-muted-foreground">Page {page + 1}</span>
        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange(page + 1)}
          disabled={!hasMore}
        >
          Next
          <ChevronRight className="w-4 h-4 ml-1" />
        </Button>
      </div>
    </div>
  );
}
