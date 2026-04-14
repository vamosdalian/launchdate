import { Button } from '@/components/ui/button';

interface ListPaginationProps {
  currentPage: number;
  totalCount: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

const ListPagination = ({ currentPage, totalCount, pageSize, onPageChange }: ListPaginationProps) => {
  const totalPages = Math.ceil(totalCount / pageSize);

  if (totalPages <= 1) {
    return null;
  }

  return (
    <div className="mt-10 flex flex-col gap-4 border-t border-[#2a2a2a] pt-6 sm:flex-row sm:items-center sm:justify-between">
      <p className="text-sm text-gray-400">
        Page {currentPage + 1} of {totalPages} · {totalCount} results
      </p>
      <div className="flex items-center gap-3">
        <Button
          variant="outline"
          className="border-[#3a3a3a] bg-[#121212] text-white hover:bg-[#1b1b1b]"
          disabled={currentPage === 0}
          onClick={() => onPageChange(currentPage - 1)}
        >
          Previous
        </Button>
        <Button
          variant="outline"
          className="border-[#3a3a3a] bg-[#121212] text-white hover:bg-[#1b1b1b]"
          disabled={currentPage >= totalPages - 1}
          onClick={() => onPageChange(currentPage + 1)}
        >
          Next
        </Button>
      </div>
    </div>
  );
};

export default ListPagination;