import itertools


# Python function to print permutations of a given list
def permutations(lst):
    # If lst is empty then there are no permutations
    if len(lst) == 0:
        return []

    # If there is only one element in lst then, only
    # one permutation is possible
    if len(lst) == 1:
        return [lst]

    # Find the permutations for lst if there are
    # more than 1 characters

    l = []  # empty list that will store current permutation

    # Iterate the input(lst) and calculate the permutation
    for i in range(len(lst)):
        m = lst[i]

        # Extract lst[i] or m from the list.  remLst is
        # remaining list
        remLst = lst[:i] + lst[i + 1:]

        # Generating all permutations where m is first
        # element
        for p in permutations(remLst):
            l.append([m] + p)
    return l


desired_colors = 3
desired_codes = 7


number_four_codes = [
    '1A2D15',
    '5D6B52',
    '94BCBD',
    '4E0104',
    '9A0606',
    'B9783F',
]
number_seven_codes = [
    'B1180B',
    'E27B52',
    'F7A711',
    'F7DBC4',
    '819B4B',
    '245F1F',
]

codes_map = {
    4: number_four_codes,
    7: number_seven_codes,
}

website_map = {
    3: 'website-1',
    4: 'website-2',
    6: 'website-3',
}

total_perms = set()
for perm_list in [permutations(list(x)) for x in itertools.combinations(codes_map[desired_codes], desired_colors)]:
    for perm in perm_list:
        total_perms.add('|'.join(perm))

for element in total_perms:
    l = element.split('|')
    joined = '-'.join(l)

    print(f'https://huemint.com/{website_map[desired_colors]}/#palette={joined}')
