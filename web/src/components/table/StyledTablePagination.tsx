import {
  Button,
  ButtonGroup,
  Flex,
  HStack,
  Select,
  Text,
  useColorModeValue as mode,
} from "@chakra-ui/react";
import React from "react";

import type { Meta } from "../../api/client/json-api-document";
import type { TablePagination } from "../../hooks/use-table-pagination";

interface IStyledTablePagination {
  my: number;
  meta: Meta;
  queryParams: TablePagination;
  setPage: (page: number) => void;
  setLimit: (limit: number) => void;
}

const StyledTablePagination = ({
  my,
  meta,
  queryParams,
  setPage,
  setLimit,
}: IStyledTablePagination) => {
  const changePage = (sign: number) => {
    if (sign < 0 && queryParams.page === 1) {
      return;
    }
    if (sign > 0 && queryParams.page === meta.totalPages) {
      return;
    }
    setPage((queryParams.page || 1) + sign);
  };

  return (
    <Flex my={my || 0} align="center" justify="space-between">
      <Text color={mode("gray.600", "gray.400")} fontSize="sm">
        Total {meta.totalItems || 0} elements | Total {meta.totalPages || 0}{" "}
        pages | Current page {queryParams.page}
      </Text>
      <HStack>
        <Select
          value={queryParams.limit}
          onChange={(event) => setLimit(+event.target.value)}
        >
          <option value="5">5</option>
          <option value="10">10</option>
          <option value="15">15</option>
          <option value="20">20</option>
        </Select>
        <ButtonGroup variant="outline">
          <Button as="a" rel="prev" onClick={() => changePage(-1)}>
            Previous
          </Button>
          <Button as="a" rel="next" onClick={() => changePage(1)}>
            Next
          </Button>
        </ButtonGroup>
      </HStack>
    </Flex>
  );
};

export default StyledTablePagination;
