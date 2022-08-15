import { Box } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

import SmersApi from "../../../api/smers-api";
import type { SmerDto } from "../../../models/smer";
import { useErrorHandler } from "../../../utils/handle-get-error";

const Smer = () => {
  const [, setSmer] = useState<SmerDto>();
  const { smerId } = useParams();
  const errorHandler = useErrorHandler();

  useEffect(() => {
    if (smerId) {
      SmersApi.view(+smerId)
        .then((payload) => {
          setSmer(payload);
        })
        .catch(errorHandler);
    }
  }, [smerId]);

  return (
    <Box as="section" py="12">
      <Box maxW={{ base: "xl", md: "7xl" }} mx="auto">
        <Box overflowX="auto" />
      </Box>
    </Box>
  );
};

export default Smer;
