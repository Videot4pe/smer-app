import { AddIcon } from "@chakra-ui/icons";
import {
  Box,
  Center,
  Heading,
  IconButton,
  Td,
  Tr,
  useDisclosure,
} from "@chakra-ui/react";
import { useState } from "react";

import SmersApi from "../../api/smers-api";
import StyledTable from "../../components/table/StyledTable";
import StyledTablePagination from "../../components/table/StyledTablePagination";
import { useTableData } from "../../hooks/use-table-data";
import { useTableFilter } from "../../hooks/use-table-filter";
import { useTablePagination } from "../../hooks/use-table-pagination";
import { useTableSort } from "../../hooks/use-table-sort";
import type { SmerDto } from "../../models/smer";
import { useErrorHandler } from "../../utils/handle-get-error";

import NewSmer from "./components/NewSmer";
import smersTableColumns from "./smers-table-columns";

const Smers = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [activeSmerId, setActiveSmerId] = useState<undefined | number>(
    undefined
  );
  const errorHandler = useErrorHandler();
  const { queryParams, setPage, setLimit } = useTablePagination();
  const { sortParams, setSortParams } = useTableSort();
  const { filterParams, arrayFilterParams, setFilterParams } = useTableFilter();
  const { data, meta, isLoading, fetch } = useTableData<SmerDto>(
    SmersApi.list,
    queryParams,
    arrayFilterParams,
    sortParams
  );

  const onRemove = (id: number) => {
    SmersApi.remove(id)
      .then(() => fetch)
      .catch(errorHandler);
  };

  const onEdit = (id: number) => {
    setActiveSmerId(id);
    onOpen();
  };

  const columns = smersTableColumns(onRemove, onEdit);

  const onSmerSave = () => {
    onClose();
    fetch();
    setActiveSmerId(undefined);
  };

  return (
    <Box as="section" py="12">
      <Box maxW={{ base: "xl", md: "7xl" }} mx="auto">
        <Box overflowX="auto">
          <Heading size="lg" mb="2">
            <div>Smers</div>
          </Heading>
          <StyledTable
            columns={columns}
            data={data}
            isLoading={isLoading}
            filterParams={filterParams}
            setSortParams={setSortParams}
            setFilterParams={setFilterParams}
          >
            <Tr>
              <Td colSpan={columns.length} p={2}>
                <Center minW="100%">
                  <IconButton
                    icon={<AddIcon />}
                    aria-label="add smer"
                    isLoading={isLoading}
                    onClick={onOpen}
                  />
                </Center>
              </Td>
            </Tr>
          </StyledTable>
          <StyledTablePagination
            my={4}
            meta={meta}
            queryParams={queryParams}
            setPage={setPage}
            setLimit={setLimit}
          />
        </Box>
      </Box>
      <NewSmer
        isOpen={isOpen}
        isLoading={isLoading}
        activeSmerId={activeSmerId}
        onClose={onClose}
        onSmerSave={onSmerSave}
      />
    </Box>
  );
};

export default Smers;
