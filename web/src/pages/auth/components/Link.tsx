import type { HTMLChakraProps } from "@chakra-ui/system";
import { chakra, useColorModeValue } from "@chakra-ui/system";

const Link = (props: HTMLChakraProps<"a">) => (
  <chakra.a
    marginStart="1"
    href="#"
    color={useColorModeValue("blue.500", "blue.200")}
    _hover={{ color: useColorModeValue("blue.600", "blue.300") }}
    display={{ base: "block", sm: "inline" }}
    {...props}
  />
);

export default Link;
