package headergov

func (h *HeaderGovModule) Unwind(num uint64) error {
	// Remove entries from h.cache that are larger than num
	h.cache.RemoveVotesAfter(num)
	h.cache.RemoveGovernanceAfter(num)

	// Update stored block numbers for votes and governance
	var voteBlockNums StoredUint64Array = h.cache.VoteBlockNums()
	WriteVoteDataBlockNums(h.ChainKv, &voteBlockNums)

	var govBlockNums StoredUint64Array = h.cache.GovBlockNums()
	WriteGovDataBlockNums(h.ChainKv, &govBlockNums)

	return nil
}
