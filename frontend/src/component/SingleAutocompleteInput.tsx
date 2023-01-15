import * as React from 'react';
import { InputAdornment, SxProps, FilterOptionsState } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import SearchIcon from '@mui/icons-material/Search';
import Autocomplete, {
  AutocompleteCloseReason,
  createFilterOptions,
} from '@mui/material/Autocomplete';
import Box from '@mui/material/Box';
import { TextField } from '@mui/material';
import { ItemType } from 'component/MultipleAutocompleteInput';

type FetcherFunc = (input: string) => Promise<ItemType[]>
type SetValueFunc = (items: ItemType) => void;

type AutocompleteInputProps = {
  disabled?: boolean
  allowAdd?: boolean
  defaultValue?: ItemType
  errorField?: any
  fetcher: FetcherFunc
  label?: string
  placeholder?: string
  setValue: SetValueFunc
  inputSx?: SxProps
  style?: React.CSSProperties
}

export default function SingleAutocompleteInput({
  disabled,
  allowAdd,
  defaultValue,
  errorField,
  fetcher,
  label,
  placeholder,
  setValue,
  inputSx,
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

  const handleOnChange = (event: React.SyntheticEvent, newItem: NonNullable<string | ItemType>, reason: string) => {
    if (event.type === 'keydown' &&
      (event as React.KeyboardEvent).key === 'Backspace' &&
      reason === 'removeOption') {
      return;
    }

    if (typeof(newItem) === 'string') {
      return;
    }

    if (newItem !== null) {
        setValue(newItem);
    }
  };

  const filter = createFilterOptions<ItemType>();

  const filterOptions = (options: ItemType[], params: FilterOptionsState<ItemType>) => {
    const filtered = filter(options, params);

    const { inputValue } = params;
    const isExisting = options.some((option) => inputValue === option.id);
    if (inputValue !== '' && !isExisting) {
      filtered.push({
        id: inputValue,
        label: inputValue,
      });
    }

    return filtered;
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
        freeSolo={allowAdd}
        selectOnFocus
        clearOnBlur
        onClose={handleOnClose}
        onChange={handleOnChange}
        disableClearable
        defaultValue={defaultValue}
        filterOptions={allowAdd ? filterOptions : undefined}
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
        getOptionLabel={(option) => typeof(option) === 'string' ? option : option.label}
        onFocus={handleOnKeyUp}
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
            sx={inputSx}
          />
        )}
      />
    </div>
  );
}
