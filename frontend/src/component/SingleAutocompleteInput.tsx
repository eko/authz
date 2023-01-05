import * as React from 'react';
import { InputAdornment, SxProps } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import SearchIcon from '@mui/icons-material/Search';
import Autocomplete, {
  AutocompleteCloseReason,
} from '@mui/material/Autocomplete';
import Box from '@mui/material/Box';
import { TextField } from '@mui/material';
import { ItemType } from 'component/MultipleAutocompleteInput';

type FetcherFunc = (input: string) => Promise<ItemType[]>
type SetValueFunc = (items: ItemType) => void;

type AutocompleteInputProps = {
  disabled?: boolean
  defaultValue?: ItemType
  errorField?: any
  fetcher: FetcherFunc
  label?: string
  placeholder?: string
  setValue: SetValueFunc
  textSx?: SxProps
  style?: React.CSSProperties
}

export default function SingleAutocompleteInput({
  disabled,
  defaultValue,
  errorField,
  fetcher,
  label,
  placeholder,
  setValue,
  textSx,
  style,
}: AutocompleteInputProps) {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const [listItems, setListItems] = React.useState<ItemType[]>([]);
  const theme = useTheme();

  const handleOnKeyUp = async (event: any) => {
    const values = await fetcher(event.target.value);
    setListItems(values);
  };

  const handleClose = () => {
    if (anchorEl) {
      anchorEl.focus();
    }
    setAnchorEl(null);
  };

  const handleOnClose = (event: React.ChangeEvent<{}>, reason: AutocompleteCloseReason) => {
    if (reason === 'escape') {
      handleClose();
    }
  };

  const handleOnChange = (event: React.SyntheticEvent, newItem: ItemType | null, reason: string) => {
    if (event.type === 'keydown' &&
      (event as React.KeyboardEvent).key === 'Backspace' &&
      reason === 'removeOption') {
      return;
    }

    if (newItem !== null) {
        setValue(newItem);
    }
  };

  React.useEffect(() => {
    if (defaultValue === undefined) {
      return;
    }

    setListItems([defaultValue]);
  }, [defaultValue]);

  return (
    <div style={style}>
      <Autocomplete
        disabled={disabled}
        openOnFocus
        onClose={handleOnClose}
        onChange={handleOnChange}
        disableClearable
        defaultValue={defaultValue}
        isOptionEqualToValue={(option: ItemType, value: ItemType): boolean => {
          return option.id === value.id;
        }}
        renderTags={() => null}
        noOptionsText='Aucun élément'
        renderOption={(props, option) => (
          <li {...props}>
            <Box
              sx={{
                flexGrow: 1,
                '& span': {
                  color: theme.palette.mode === 'light' ? '#586069' : '#8b949e',
                },
              }}
            >
              {option.label}
            </Box>
          </li>
        )}
        options={[...listItems]}
        getOptionLabel={(option) => option.label}
        onKeyUp={handleOnKeyUp}
        renderInput={(params) => (
          <TextField
            ref={params.InputProps.ref}
            disabled={disabled}
            role='singleautocompleteinput-field'
            inputProps={params.inputProps}
            label={label}
            placeholder={placeholder}
            error={errorField ? true : false}
            helperText={errorField?.message}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon />
                </InputAdornment>
              ),
            }}
            sx={textSx}
          />
        )}
      />
    </div>
  );
}
