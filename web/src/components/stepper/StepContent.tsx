import type { CenterProps } from "@chakra-ui/react";
import { AbsoluteCenter, Center } from "@chakra-ui/react";

const StepContent = (props: CenterProps) => (
  <AbsoluteCenter height="full" width="full" position="absolute" zIndex={1}>
    <Center height="full" fontSize="sm" fontWeight="semibold" {...props} />
  </AbsoluteCenter>
);

export default StepContent;
