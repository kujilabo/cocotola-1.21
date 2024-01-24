import { FC } from 'react';

import { Button } from '@chakra-ui/react';

type BasicButtonProps = {
  value: string;
  isDisabled?: boolean;
  isLoading?: boolean;
  onClick?: () => void;
};

export const BasicButton: FC<BasicButtonProps> = (props: BasicButtonProps) => {
  return (
    <Button
      colorScheme="teal"
      isDisabled={props.isDisabled}
      isLoading={props.isLoading}
      onClick={props.onClick}
    >
      {props.value}
    </Button>
  );
};
