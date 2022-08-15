import { Center } from "@chakra-ui/react";
import { Audio } from "react-loader-spinner";

const Loader: React.FC = () => {
  return (
    <Center h="100vh" w="100vw">
      <Audio height="100" width="100" color="grey" ariaLabel="loading" />
    </Center>
  );
};

export default Loader;
