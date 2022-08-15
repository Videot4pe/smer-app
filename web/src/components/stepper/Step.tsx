import { Box, Flex, useColorModeValue } from "@chakra-ui/react";
import type * as React from "react";

import SkewBox from "./SkewBox";
import StepContent from "./StepContent";

interface StepProps {
  isCurrent?: boolean;
  isDone?: boolean;
  disabled?: boolean;
  onClick?: () => void;
  children: React.ReactNode;
}

const Step = (props: StepProps) => {
  const { children, isCurrent, disabled, isDone, onClick } = props;
  const color = useColorModeValue("white", "gray.900");
  return (
    <Box as="li" flex="1">
      <Box
        as="button"
        outline={0}
        className="group"
        width="full"
        disabled={disabled}
        onClick={onClick}
      >
        <Flex
          align="center"
          height="12"
          justify="center"
          position="relative"
          css={{ "--arrow-skew": "20deg" }}
        >
          <SkewBox isCurrent={isCurrent} isDone={isDone} placement="top" />
          <StepContent color={isCurrent ? color : "inherit"}>
            {children}
          </StepContent>
          <SkewBox isCurrent={isCurrent} isDone={isDone} placement="bottom" />
        </Flex>
      </Box>
    </Box>
  );
};

export default Step;
