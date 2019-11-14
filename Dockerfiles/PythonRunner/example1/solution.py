def two_oldest_ages(ages):
    oldest = 0
    next_oldest = 0
    for age in ages:
        if age > oldest:
            next_oldest = oldest
            oldest = age
        elif age > next_oldest:
            next_oldest = age
    return [next_oldest, oldest]